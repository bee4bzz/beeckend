from typing import TYPE_CHECKING

from sqlalchemy import Column, String
from fastapi_users_db_sqlalchemy import SQLAlchemyBaseUserTableUUID

from app.db.base_class import Base, Base1toN


class User(
    SQLAlchemyBaseUserTableUUID,
    Base1toN,
    Base,  # type: ignore
):
    """
    User model.
    """

    __ownedtablename__ = "Cheptel"
    full_name = Column(String, nullable=False, index=True)
