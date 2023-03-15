from fastapi_users_db_sqlalchemy import SQLAlchemyUserDatabase
    
class CRUDUser(SQLAlchemyUserDatabase):
    """
    CRUD for :any:`UserModel`.
    """
