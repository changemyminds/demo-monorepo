from fastapi.testclient import TestClient
from service_c.main import app

client = TestClient(app)


def test_healthz():
    resp = client.get("/healthz")
    assert resp.status_code == 200
    assert resp.json()["status"] == "ok"


def test_version():
    resp = client.get("/version")
    assert resp.status_code == 200
    assert resp.json()["service"] == "service-c"
    assert "version" in resp.json()


def test_root_uses_shared_lib():
    resp = client.get("/", params={"name": "Ada"})
    assert resp.status_code == 200
    assert "Ada" in resp.json()["message"]
