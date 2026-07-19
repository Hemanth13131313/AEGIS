package aegis.policy.default

default allow = true

allow {
    input.method == "GET"
}
