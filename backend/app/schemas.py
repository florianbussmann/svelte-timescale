# schemas.py
from pydantic import BaseModel
from datetime import datetime


class SensorDataBase(BaseModel):
    device_id: str
    value: int


class SensorDataCreate(SensorDataBase):
    pass


class SensorDataRead(SensorDataBase):
    timestamp: datetime

    class Config:
        from_attributes = True
