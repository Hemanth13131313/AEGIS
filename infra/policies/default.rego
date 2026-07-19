package aegis.policy.core

default allow = true

allow {
    input.method == "GET"
}
