# from sqlalchemy.orm import Session
from sqlalchemy.ext.asyncio import AsyncSession
from app.core.security import password_helper
from app import crud, models, schemas
from app.core.config import settings
from app import db # noqa: F401

# make sure all SQL Alchemy models are imported (app.db.base) before initializing DB
# otherwise, SQL Alchemy might fail to initialize relationships properly
# for more details: https://github.com/tiangolo/full-stack-fastapi-postgresql/issues/28


async def init_db(db: AsyncSession) -> None:
    """
    Initialize the database with some data.
    """
    # Tables should be created with Alembic migrations
    # But if you don't want to use migrations, create
    # the tables un-commenting the next line
    # Base.metadata.create_all(bind=engine)

    user_crud = crud.CRUDUser(db, models.User)

    user = await user_crud.get_by_email(email=settings.FIRST_SUPERUSER)
    if not user:
        user_in = schemas.UserCreate(
            email=settings.FIRST_SUPERUSER,
            password=settings.FIRST_SUPERUSER_PASSWORD,
            full_name=settings.FIRST_SUPERUSER_NAME,
            is_superuser=True,
            is_verified=True,
        ).dict()
        user_in['hashed_password'] = password_helper.hash(user_in.pop('password'))
        user = await user_crud.create(create_dict=user_in)  # noqa: F841
        user = await user_crud.update(
            user=user, update_dict={"is_active": True}
        )  # noqa: F841
