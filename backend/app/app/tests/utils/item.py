from typing import Optional
import uuid

# from sqlalchemy.orm import Session
from sqlalchemy.ext.asyncio import AsyncSession

from app import crud, models
from app.models.item import Item
from app.schemas.item import ItemCreate
from app.tests.utils.user import create_random_user
from app.tests.utils.utils import random_lower_string


async def create_random_item(
    db: AsyncSession, *, owner_id: Optional[uuid.UUID] = None
) -> models.Item:
    item_crud = crud.CRUDItem(db, Item)
    if owner_id is None:
        user = await create_random_user(db)
        owner_id = user.id
    title = random_lower_string()
    description = random_lower_string()
    item_in = ItemCreate(title=title, description=description)
    return await item_crud.create_with_owner(create_dict=item_in, owner_id=owner_id)
