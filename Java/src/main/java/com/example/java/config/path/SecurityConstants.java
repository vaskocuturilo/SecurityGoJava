package com.example.java.config.path;

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

    public static final String[] ROLE_API_CLIENT = {
            "USER",
    };

}