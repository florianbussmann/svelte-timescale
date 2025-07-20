from typing import List
from fastapi import FastAPI, Depends
from fastapi.middleware.cors import CORSMiddleware

from sqlalchemy.orm import Session

from models import Base, SensorData
from db import engine, get_db
from schemas import SensorDataCreate, SensorDataRead

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
