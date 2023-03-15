from fastapi import APIRouter
from app.api.api_v1.users import users

from app.schemas import UserRead, UserCreate, UserUpdate


api_router_user: APIRouter = APIRouter()

api_router_user.include_router(
    users.users_manager.get_auth_router(users.auth_backend, requires_verification=True),
    prefix="/auth/jwt",
    tags=["auth"],
)
api_router_user.include_router(
    users.users_manager.get_register_router(UserRead, UserCreate),
    prefix="/auth",
    tags=["auth"],
)
api_router_user.include_router(
    users.users_manager.get_reset_password_router(),
    prefix="/auth",
    tags=["auth"],
)
api_router_user.include_router(
    users.users_manager.get_verify_router(UserRead),
    prefix="/auth",
    tags=["auth"],
)
api_router_user.include_router(
    users.users_manager.get_users_router(UserRead, UserUpdate, requires_verification=True),
    prefix="/users",
    tags=["users"],
)
