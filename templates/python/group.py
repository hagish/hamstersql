from sqlalchemy import text

import core

{{range .}}
# {{.ReturnTypeStr }}: {{.Doc}}
{{- if .Parameters}}
def {{.Properties.name}}(session{{range .Parameters}}, {{.}}{{end}}):
{{- else}}
def {{.Properties.name}}(session):
{{- end}}
    query = text("""
{{.SqlQuery}}
    """)
    params = {
{{- range .Parameters}}
        "{{.}}": {{.}},
{{- end}}
    }
    result = session.execute(query, params)
    {{- if eq .ReturnTypeStr "OneRow"}}
    return core.convert_result_to_one_row(result)
    {{- end}}
    {{- if eq .ReturnTypeStr "ManyRows"}}
    return core.convert_result_to_many_rows(result)
    {{- end}}
    {{- if eq .ReturnTypeStr "AffectedRows"}}
    return core.convert_result_to_affected_rows(result)
    {{- end}}
    {{- if eq .ReturnTypeStr "Scalar"}}
    return core.convert_result_to_scalar(result)
    {{- end}}
    {{- if eq .ReturnTypeStr "InsertID"}}
    return core.convert_result_to_insert_id(result)
    {{- end}}
    {{- if eq .ReturnTypeStr "None"}}
    return result
    {{- end}}

{{end}}
