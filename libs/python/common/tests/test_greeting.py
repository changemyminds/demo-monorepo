from common import greeting


def test_greeting_with_name():
    assert "Ada" in greeting("Ada")


def test_greeting_default():
    assert "world" in greeting()
