from typing import TYPE_CHECKING
import uuid

from sqlalchemy import Column, String
from app.db.base_class import ID, Base1toN

if TYPE_CHECKING:
    from .user import User  # noqa: F401


class ItemBase:
    """
    Base class for item models.
    """

    id: Column[ID] = Column(ID, primary_key=True, default=uuid.uuid4, nullable=True)
    title: Column[String] = Column(String, index=True, nullable=False)
    type: Column[String] = Column(String, index=True, nullable=False)
