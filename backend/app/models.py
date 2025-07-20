from sqlalchemy import Column, Integer, Text, TIMESTAMP, func
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()


class SensorData(Base):
    __tablename__ = "sensor_data"

    device_id = Column(Text, nullable=False)
    timestamp = Column(
        TIMESTAMP(timezone=True), primary_key=True, server_default=func.now()
    )
    value = Column(Integer)
