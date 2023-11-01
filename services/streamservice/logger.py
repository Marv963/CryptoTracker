import logging
from logging.handlers import TimedRotatingFileHandler
import os

from .config import Config


class Logger:
    logger: logging.Logger

    @classmethod
    def initialize_logger(
        cls, log_file_name: str, log_level: int = Config.LOGGER
    ) -> None:
        cls.logger = logging.getLogger(log_file_name)
        cls.logger.setLevel(log_level)

        file_handler = TimedRotatingFileHandler(
            os.path.join("logs", log_file_name),
            when="midnight",
            interval=1,
            backupCount=28,
        )
        file_handler.setLevel(log_level)
        file_handler.setFormatter(CustomFormatter())
        cls.logger.addHandler(file_handler)

        console_handler = logging.StreamHandler()
        console_handler.setLevel(log_level)
        console_handler.setFormatter(CustomFormatter())
        cls.logger.addHandler(console_handler)

    @classmethod
    def info(cls, message: str) -> None:
        cls.logger.info(message)

    @classmethod
    def debug(cls, message: str) -> None:
        cls.logger.debug(message)

    @classmethod
    def warning(cls, message: str) -> None:
        cls.logger.warning(message)

    @classmethod
    def error(cls, message: str) -> None:
        cls.logger.error(message)

    @classmethod
    def critical(cls, message: str) -> None:
        cls.logger.critical(message)


class CustomFormatter(logging.Formatter):
    grey = "\x1b[38;20m"
    yellow = "\x1b[33;20m"
    red = "\x1b[31;20m"
    green = "\x1b[1;32m"
    bold_red = "\x1b[31;1m"
    reset = "\x1b[0m"
    log_format = (
        "%(asctime)s - %(name)s - %(levelname)s - %(message)s (%(filename)s:%(lineno)d)"
    )

    FORMATS = {
        logging.DEBUG: green + log_format + reset,
        logging.INFO: grey + log_format + reset,
        logging.WARNING: yellow + log_format + reset,
        logging.ERROR: red + log_format + reset,
        logging.CRITICAL: bold_red + log_format + reset,
    }

    def format(self, record: logging.LogRecord) -> str:
        log_fmt = self.FORMATS.get(record.levelno)
        formatter = logging.Formatter(log_fmt)
        return formatter.format(record)
