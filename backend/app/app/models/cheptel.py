from typing import TYPE_CHECKING

from app.models import ItemBase
from app.db.base_class import Base1toN, BaseNto1, Base

if TYPE_CHECKING:
    from .user import User  # noqa: F401


class Cheptel(
    Base,  # type: ignore
    BaseNto1,
    Base1toN,
    ItemBase,
):
    __ownertablename__ = "User"
    __ownedtablename__ = "Hive"
