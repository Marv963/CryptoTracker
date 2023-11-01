import os
import logging


class Config:
    ##### GENERAL #####
    APP = os.environ.get("APP", "DEV")

    ##### LOGGING #####
    LOGGING: str = os.environ.get("LOGGING", "WARNING")
    LOGGER: int = (
        logging.DEBUG
        if LOGGING == "DEBUG"
        else logging.INFO
        if LOGGING == "INFO"
        else logging.WARNING
        if LOGGING == "WARNING"
        else logging.ERROR
        if LOGGING == "ERROR"
        else logging.CRITICAL
        if LOGGING == "CRITICAL"
        else logging.WARNING
    )

    ##### RABBITMQ #####
    RABBITMQ_USERNAME: str = os.environ.get("RABBITMQ_USERNAME", "")
    RABBITMQ_PASSWORD: str = os.environ.get("RABBITMQ_PASSWORD", "")
    RABBITMQ_HOST: str = os.environ.get("RABBITMQ_HOST", "")
    RABBITMQ_PORT: int = int(os.environ.get("RABBITMQ_PORT", 0))
    RABBITMQ_EXCHANGE: str = os.environ.get("RABBITMQ_EXCHANGE", "")
    RABBITMQ_QUEUE: str
    RABBITMQ_ROUTING_KEY: str

    @classmethod
    def validate(cls) -> None:
        missing_vars = [
            var for var, value in cls.__dict__.items() if value == "" or value == 0
        ]
        if missing_vars:
            logging.error(f"Missing environment variables: {', '.join(missing_vars)}")
            raise ValueError("Missing configuration")
