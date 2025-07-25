from typing import List
from fastapi import FastAPI, Depends
from fastapi.middleware.cors import CORSMiddleware

from sqlalchemy.orm import Session
from sqlalchemy.dialects.postgresql import insert

from models import Base, Sensor, SensorData
from db import engine, get_db
from schemas import (
    SensorDataCreate,
    SensorDataRead,
    SensorRead,
    SensorUpsert,
)

app = FastAPI()


origins = [
    "*",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/")
def read_root():
    return {"Hello": "World"}


@app.post("/sensors/", response_model=SensorRead)
def upsert_sensor(sensor: SensorUpsert, db: Session = Depends(get_db)):
    # Convert to dict, skip unset fields
    update_fields = sensor.dict(exclude_unset=True, exclude={"id"})

    # If no fields to update, skip update part
    stmt = insert(Sensor).values(**sensor.dict())

    if update_fields:
        stmt = stmt.on_conflict_do_update(index_elements=["id"], set_=update_fields)
    else:
        # do nothing on conflict if no fields to update
        stmt = stmt.on_conflict_do_nothing(index_elements=["id"])

    db.execute(stmt)
    db.commit()

    return db.query(Sensor).filter(Sensor.id == sensor.id).first()


@app.get("/sensors/", response_model=List[SensorRead])
def read_sensors(db: Session = Depends(get_db)):
    return db.query(Sensor).all()


@app.get("/data/", response_model=List[SensorDataRead])
def read_data(skip: int = 0, limit: int = 100, db: Session = Depends(get_db)):
    return db.query(SensorData).offset(skip).limit(limit).all()


@app.post("/data/", response_model=SensorDataRead)
def create_data(data: SensorDataCreate, db: Session = Depends(get_db)):
    db_item = SensorData(**data.dict())
    db.add(db_item)
    db.commit()
    db.refresh(db_item)
    return db_item
