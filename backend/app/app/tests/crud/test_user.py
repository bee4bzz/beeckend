import pytest

# from sqlalchemy.orm import Session
from sqlalchemy.ext.asyncio import AsyncSession

from app import crud, models
from app.core.security import password_helper
from app.schemas.user import UserCreate, UserUpdate
from app.tests.utils.utils import random_email, random_lower_string

pytestmark = pytest.mark.asyncio


async def test_create_user(async_get_db: AsyncSession) -> None:
    user_crud = crud.CRUDUser(async_get_db, models.User)
    email = random_email()
    password = random_lower_string()
    user_in = UserCreate(email=email, password=password)
    user = await user_crud.create(create_dict=user_in.dict())
    assert user.email == email
    assert hasattr(user, "hashed_password")


async def test_check_if_user_is_active(async_get_db: AsyncSession) -> None:
    user_crud = crud.CRUDUser(async_get_db, models.User)
    email = random_email()
    password = random_lower_string()
    user_in = UserCreate(email=email, password=password)
    user = await user_crud.create(create_dict=user_in.dict())
    is_active = user.active
    assert is_active is True


async def test_check_if_user_is_active_inactive(async_get_db: AsyncSession) -> None:
    user_crud = crud.CRUDUser(async_get_db, models.User)
    email = random_email()
    password = random_lower_string()
    user_in = UserCreate(email=email, password=password, is_active=False)
    user = await user_crud.create(create_dict=user_in.dict())
    is_active = user.active
    assert is_active


async def test_check_if_user_is_superuser(async_get_db: AsyncSession) -> None:
    user_crud = crud.CRUDUser(async_get_db, models.User)
    email = random_email()
    password = random_lower_string()
    user_in = UserCreate(email=email, password=password, is_superuser=True)
    user = await user_crud.create(create_dict=user_in.dict())
    is_superuser = user.superuser
    assert is_superuser is True


async def test_check_if_user_is_superuser_normal_user(
    async_get_db: AsyncSession,
) -> None:
    user_crud = crud.CRUDUser(async_get_db, models.User)
    email = random_email()
    password = random_lower_string()
    user_in = UserCreate(email=email, password=password, is_superuser=True)
    user = await user_crud.create(create_dict=user_in.dict())
    is_superuser = user.superuser
    assert is_superuser is False


async def test_get_user(async_get_db: AsyncSession) -> None:
    user_crud = crud.CRUDUser(async_get_db, models.User)
    password = random_lower_string()
    username = random_email()
    user_in = UserCreate(email=username, password=password, is_superuser=True)
    user = await user_crud.create(create_dict=user_in.dict())
    user_2 = await user_crud.get(id=user.id)
    assert user_2
    assert user.email == user_2.email
    assert user.dict() == user_2.dict()


async def test_update_user(async_get_db: AsyncSession) -> None:
    user_crud = crud.CRUDUser(async_get_db, models.User)
    password = random_lower_string()
    email = random_email()
    user_in = UserCreate(email=email, password=password, is_superuser=True)
    user = await user_crud.create(create_dict=user_in.dict())
    new_password = random_lower_string()
    user_in_update = UserUpdate(password=new_password, is_superuser=True)
    await user_crud.update(user=user, update_dict=user_in_update.dict())
    user_2 = await user_crud.get(id=user.id)
    assert user_2
    assert user.email == user_2.email
    assert password_helper.verify_and_update(new_password, user_2.hashed_password)
