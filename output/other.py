# generated by hamstersql
from sqlalchemy import text

import core


# one: 
def get_something(session):
    query = text("""
select * from others limit 1
    """)
    params = {
    }
    result = session.execute(query, params)
    return core.convert_result_to_one_row(result)


