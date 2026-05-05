package com.example.java.dto;

import java.util.Arrays;


public record CredentialsDto(String email, char[] password) {
    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;

        CredentialsDto that = (CredentialsDto) o;

        if (!email.equals(that.email)) return false;
        return Arrays.equals(password, that.password);
    }

    @Override
    public int hashCode() {
        int result = email.hashCode();
        result = 31 * result + Arrays.hashCode(password);
        return result;
    }

    @Override
    public String toString() {
        return "CredentialsDto{" +
                "email='" + email + '\'' +
                ", password=" + Arrays.toString(password) +
                '}';
    }
}
