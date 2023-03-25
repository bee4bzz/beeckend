from typing import TYPE_CHECKING
from app.models import ItemBase
from app.db.base_class import BaseNto1, Base

if TYPE_CHECKING:
    from .cheptel import Cheptel  # noqa: F401


class Hive(
    Base,  # type: ignore
    BaseNto1,
    ItemBase,
):
    __ownertablename__ = "Cheptel"
