package com.example.java.dto;

import java.util.Arrays;
import java.util.Objects;

public record SignUpDto(String email, char[] password) {
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
                ", password=" + Arrays.toString(password) +
                '}';
    }
}