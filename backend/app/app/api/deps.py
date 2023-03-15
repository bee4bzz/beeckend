from typing import Generator, AsyncGenerator

from fastapi import Depends
from sqlalchemy.ext.asyncio import AsyncSession

from app import crud, models
from app.db.session import SessionLocal
from app.db.session import async_session


def get_db() -> Generator:
    """
    Dependency to get a db session.
    """
    try:
        db = SessionLocal()
        yield db
    finally:
        db.close()


async def async_get_db() -> AsyncGenerator:
    """
    Dependency to get an async db session.
    """
    async with async_session() as session:
        yield session


async def get_user_db(
    session: AsyncSession = Depends(async_get_db),
) -> AsyncGenerator[crud.CRUDUser, None]:
    """
    Dependency to get a user crud.
    """
    yield crud.CRUDUser(session, models.User)
