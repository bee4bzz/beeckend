import uuid
from typing import Optional
from fastapi_users import schemas
from pydantic import BaseModel


# Shared properties
class UserBase(BaseModel):
    full_name: Optional[str] = None


class UserRead(schemas.BaseUser[uuid.UUID], UserBase):
    pass


class UserCreate(schemas.BaseUserCreate, UserBase):
    pass


class UserUpdate(schemas.BaseUserUpdate, UserBase):
    pass
