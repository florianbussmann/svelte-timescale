# schemas.py
from typing import Optional
from pydantic import BaseModel
from datetime import datetime


class SensorBase(BaseModel):
    name: Optional[str] = None


class SensorUpsert(SensorBase):
    id: str  # required for identifying the sensor


class SensorRead(SensorBase):
    id: str

    class Config:
        from_attributes = True


class SensorDataBase(BaseModel):
    device_id: str
    value: int


class SensorDataCreate(SensorDataBase):
    pass


class SensorDataRead(SensorDataBase):
    timestamp: datetime

    class Config:
        from_attributes = True
