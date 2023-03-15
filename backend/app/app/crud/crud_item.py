from typing import List
import uuid

# from sqlalchemy.orm import Session
from sqlalchemy.sql.expression import select

from app.crud.base import CRUDBase1ToN
from app.models.item import ItemBase
from app.schemas.item import ItemCreate, ItemUpdate


class CRUDItem(CRUDBase1ToN[ItemBase, ItemCreate, ItemUpdate]):
    """CRUD for :any:`Item model<models.item.Item>"""

    async def create_with_owner(
        self, *, create_dict: ItemCreate, owner_id: uuid.UUID
    ) -> ItemBase:
        db_obj = self.model(**create_dict, owner_id=owner_id)
        self.db.add(db_obj)
        await self.db.commit()
        await self.db.refresh(db_obj)
        return db_obj

    async def get_multi_by_owner(
        self, *, owner_id: int, skip: int = 0, limit: int = 100
    ) -> List[ItemBase]:
        result = await self.db.execute(
            select(self.model)
            .filter(ItemBase.owner_id == owner_id)
            .offset(skip)
            .limit(limit)
            .all()
        )
        return result.scalars().all()
