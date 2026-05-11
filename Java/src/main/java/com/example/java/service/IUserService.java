package com.example.java.service;

import com.example.java.dto.*;

public interface IUserService {

    AuthResponse login(final CredentialsDto credentials);

    AuthResponse register(final SignUpDto userDto);

    AuthResponse refresh(final RefreshTokenRequest request);


    void logoutAll(final String email);
}
