from typing import Any, Dict, Generic, List, Optional, Type, TypeVar, Union, cast
from fastapi.encoders import jsonable_encoder
from pydantic import BaseModel
from sqlalchemy.ext.asyncio import AsyncSession

from app.db.base_class import ID, Base, Base1toN, BaseNto1
from sqlalchemy.sql.expression import select

ModelType = TypeVar("ModelType", bound=Base)
CreateSchemaType = TypeVar("CreateSchemaType", bound=BaseModel)
UpdateSchemaType = TypeVar("UpdateSchemaType", bound=BaseModel)


class CRUDBase(Generic[ModelType, CreateSchemaType, UpdateSchemaType]):
    def __init__(self, db: AsyncSession, model: Type[ModelType]):
        """
        CRUD object with default methods to Create, Read, Update, Delete (CRUD).

        **Parameters**

        * `model`: A SQLAlchemy model class
        * `schema`: A Pydantic model (schema) class
        """
        self.db = db
        self.model = model

    async def get(self, id: Any) -> Optional[ModelType]:
        result = await self.db.execute(select(self.model).filter(self.model.id == id))
        return result.scalars().first()

    async def get_multi(self, *, skip: int = 0, limit: int = 100) -> List[ModelType]:
        query = select(self.model).offset(skip).limit(limit)
        result = await self.db.execute(query)
        res = result.scalars().all()
        return res

    async def create(self, *, create_dict: CreateSchemaType) -> ModelType:
        create_dict_data = create_dict.dict()
        db_obj = self.model(**create_dict_data)  # type: ignore
        self.db.add(db_obj)
        await self.db.commit()
        await self.db.refresh(db_obj)
        return db_obj

    async def update(
        self, *, db_obj: ModelType, create_dict: Union[UpdateSchemaType, Dict[str, Any]]
    ) -> ModelType:
        obj_data = jsonable_encoder(db_obj)
        if isinstance(create_dict, dict):
            update_data = create_dict
        else:
            update_data = create_dict.dict(exclude_unset=True)
        for field in obj_data:
            if field in update_data:
                setattr(db_obj, field, update_data[field])
        self.db.add(db_obj)
        await self.db.commit()
        await self.db.refresh(db_obj)
        return db_obj

    async def remove(self, *, id: ID) -> ModelType:
        obj = await self.db.get(self.model, id)
        assert obj is not None
        await self.db.delete(obj)
        await self.db.commit()
        return cast(ModelType, obj)


ModelType = TypeVar("ModelType", bound=Base1toN)


class CRUDBase1ToN(
    Generic[ModelType, CreateSchemaType, UpdateSchemaType],
):
    def __init__(self, db: AsyncSession, model: Type[ModelType]):
        """
        CRUD object with default methods to Create, Read, Update, Delete (CRUD).

        **Parameters**

        * `model`: A SQLAlchemy model class
        * `schema`: A Pydantic model (schema) class
        """
        self.db = db
        self.model = model

    async def get_children(
        self, *, owner_id: int, skip: int = 0, limit: int = 100
    ) -> List[ModelType]:
        result = await self.db.execute(
            select(self.model)
            .filter(self.model.id == owner_id)
            .offset(skip)
            .limit(limit)
        )
        return result.scalars().all()


ModelType = TypeVar("ModelType", bound=BaseNto1)


class CRUDBaseNTo1(Generic[ModelType, CreateSchemaType, UpdateSchemaType]):
    def __init__(self, db: AsyncSession, model: Type[ModelType]):
        """
        CRUD object with default methods to Create, Read, Update, Delete (CRUD).

        **Parameters**

        * `model`: A SQLAlchemy model class
        * `schema`: A Pydantic model (schema) class
        """
        self.db = db
        self.model = model

    async def get_multi_by_owner(
        self, *, owner_id: int, skip: int = 0, limit: int = 100
    ) -> List[ModelType]:
        result = await self.db.execute(
            select(self.model)
            .filter(self.model.owner_id == owner_id)
            .offset(skip)
            .limit(limit)
        )
        return result.scalars().all()
