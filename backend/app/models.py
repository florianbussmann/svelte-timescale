from sqlalchemy import Column, Integer, Text, TIMESTAMP, func, ForeignKey
from sqlalchemy.orm import relationship
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()


class Sensor(Base):
    __tablename__ = "sensors"

    id = Column(Text, primary_key=True)
    name = Column(Text)

    data = relationship("SensorData", back_populates="sensor")


class SensorData(Base):
    __tablename__ = "sensor_data"

    device_id = Column(Text, ForeignKey("sensors.id"), nullable=False)
    timestamp = Column(
        TIMESTAMP(timezone=True), primary_key=True, server_default=func.now()
    )
    value = Column(Integer)

    sensor = relationship("Sensor", back_populates="data")
