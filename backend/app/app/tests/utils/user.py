from typing import Dict

# from fastapi.testclient import TestClient
from httpx import AsyncClient

from pydantic import EmailStr

# from sqlalchemy.orm import Session
from sqlalchemy.ext.asyncio import AsyncSession

from app import crud
from app.core.config import settings
from app.models.user import User
from app.schemas.user import UserCreate, UserUpdate
from app.tests.utils.utils import random_email, random_lower_string


async def user_authentication_headers(
    *, client: AsyncClient, email: str, password: str
) -> Dict[str, str]:
    data = {"email": email, "password": password}

    r = await client.post(f"{settings.API_V1_STR}/login/access-token", data=data)
    response = r.json()
    auth_token = response["access_token"]
    headers = {"Authorization": f"Bearer {auth_token}"}
    return headers


async def create_random_user(db: AsyncSession) -> User:
    user_crud = crud.CRUDUser(db, User)
    email = random_email()
    password = random_lower_string()
    user_in = UserCreate(full_name=email, email=email, password=password)
    user = await user_crud.create(create_dict=user_in.dict())
    return user


async def authentication_token_from_email(
    *, client: AsyncClient, email: EmailStr, db: AsyncSession
) -> Dict[str, str]:
    """
    Return a valid token for the user with given email.

    If the user doesn't exist it is created first.
    """
    user_crud = crud.CRUDUser(db, User)
    password = random_lower_string()
    user = await user_crud.get_by_email(email=email)
    if not user:
        user_in_create = UserCreate(full_name=email, email=email, password=password)
        user = await user_crud.create(create_dict=user_in_create.dict())
    else:
        user_in_update = UserUpdate(password=password, **user.dict())
        user = await user_crud.update(user=user, update_dict=user_in_update.dict())

    return await user_authentication_headers(
        client=client, email=email, password=password
    )
