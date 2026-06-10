package com.example.java.config.path;

import java.util.Set;

public final class SecurityConstants {
    private SecurityConstants() {
        /* This utility class should not be instantiated */
    }

    public static final String[] PUBLIC_ROUTES = {
            "/api/v1/users/login",
            "/api/v1/users/register",
            "/api/v1/users/logout",
            "/api/v1/users/refresh"

    };

    public static final String PRIVATE_ROUTE = "/api/v1/tasks/**";

    public static final Set<String> ALLOWED_ALGORITHMS = Set.of("HS256");

}