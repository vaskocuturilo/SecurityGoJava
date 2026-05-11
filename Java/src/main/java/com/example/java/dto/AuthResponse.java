package com.example.java.dto;

public record AuthResponse(
        UserDto user,
        String accessToken,
        String refreshToken) {
}