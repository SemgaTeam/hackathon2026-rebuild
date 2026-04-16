#!/bin/sh
NODE_ID="$$(docker compose exec test-garage /garage -c /etc/garage.toml node id | tail -n1 | cut -d@ -f1)"
TEST_ACCESS_KEY_ID=GK1834094781786f8dde242381
TEST_SECRET_KEY=6cb5fe16ca3df92f3c6700de488fd90d4b84802a6e89e5da7445a9274d23765d
if ! docker compose exec test-garage /garage -c /etc/garage.toml layout show | grep -q "$$NODE_ID"; then
  docker compose exec test-garage /garage -c /etc/garage.toml layout assign -z test -c 1G "$$NODE_ID"
  docker compose exec test-garage /garage -c /etc/garage.toml layout apply --version 1
fi
if ! docker compose exec test-garage /garage -c /etc/garage.toml key list | grep -q "test_access_key"; then
  docker compose exec test-garage /garage -c /etc/garage.toml key import --yes -n test_access_key $TEST_ACCESS_KEY_ID $TEST_SECRET_KEY
fi
if ! docker compose exec test-garage /garage -c /etc/garage.toml bucket list | grep -q "semga-test"; then
  docker compose exec test-garage /garage -c /etc/garage.toml bucket create semga-test
fi
docker compose exec test-garage /garage -c /etc/garage.toml bucket allow --read --write --owner --key test_access_key semga-test
