from typing import TYPE_CHECKING

from sqlalchemy import Column, String
from fastapi_users_db_sqlalchemy import SQLAlchemyBaseUserTableUUID

from app.db.base_class import Base, Base1toN

if TYPE_CHECKING:
    from .item import ItemBase  # noqa: F401


class User(
    SQLAlchemyBaseUserTableUUID,
    Base1toN,
    Base,  # type: ignore
):
    """
    User model.
    """

    __ownedtablename__ = "User"
    full_name = Column(String, nullable=False, index=True)
