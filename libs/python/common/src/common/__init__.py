"""Shared helpers for the Python services.

Consumed live-at-head via the uv workspace — no version pin, no publish.
A change here instantly affects every consuming service.
"""


def greeting(name: str = "world") -> str:
    """Build the standard greeting used by all Python services."""
    if not name:
        name = "world"
    return f"Hello, {name}! (from libs/python/common)"


__all__ = ["greeting"]
