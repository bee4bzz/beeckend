from sqlalchemy import Column, ForeignKey
from sqlalchemy.ext.declarative import declared_attr, declarative_base
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import RelationshipProperty, relationship

ID = UUID


class BaseModel:
    id: ID
    __name__: str

    # Generate __tablename__ automatically
    @declared_attr
    def __tablename__(cls) -> str:
        return cls.__name__.lower()


Base = declarative_base(cls=BaseModel)


class BaseNto1(BaseModel):
    """
    Base class for models with a N-to-1 relationship to another model.
    """

    __ownertablename__: str

    @declared_attr
    def owner_id(cls) -> Column[ID]:
        return Column(
            ID, ForeignKey(f"{cls.__ownertablename__.lower()}.id"), nullable=False
        )

    @declared_attr
    def owner(cls) -> RelationshipProperty:
        return relationship(
            cls.__ownertablename__, back_populates=f"{cls.__name__.lower()}s"
        )


class Base1toN(BaseModel):
    """
    Base class for models with a 1-to-N relationship to a model.
    """

    __ownedtablename__: str

    @declared_attr
    def children(cls) -> RelationshipProperty:
        return relationship(cls.__ownedtablename__, back_populates="owner")
