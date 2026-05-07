package com.example.java.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;

import java.util.Arrays;
import java.util.Objects;

public record SignUpDto(
        @Email(message = "Email is not valid", regexp = "^[a-zA-Z0-9_!#$%&'*+/=?`{|}~^.-]+@[a-zA-Z0-9.-]+$")
        @NotEmpty(message = "Email cannot be empty") String email,

        @NotNull(message = "Password cannot be null")
        @Size(min = 8, message = "Password must be at least 8 characters")
        char[] password) {
    @Override
    public boolean equals(Object o) {
        if (!(o instanceof SignUpDto(String email1, char[] password1))) return false;

        return Objects.equals(email(), email1) && Arrays.equals(password(), password1);
    }

    @Override
    public int hashCode() {
        int result = Objects.hashCode(email());
        result = 31 * result + Arrays.hashCode(password());
        return result;
    }

    @Override
    public String toString() {
        return "SignUpDto{" +
                "email='" + email + '\'' +
                ", password=[PROTECTED]" +
                '}';
    }
}