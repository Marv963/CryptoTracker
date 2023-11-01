from typing import Any, Dict


class TickerInfo(Dict[str, Any]):
    pass


class Ticker(Dict[str, TickerInfo]):
    pass
