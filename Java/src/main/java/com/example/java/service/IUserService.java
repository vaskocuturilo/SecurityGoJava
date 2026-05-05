package com.example.java.service;

import com.example.java.dto.CredentialsDto;
import com.example.java.dto.SignUpDto;
import com.example.java.dto.UserDto;

public interface IUserService {

    UserDto login(final CredentialsDto credentials);

    UserDto register(final SignUpDto userDto);
}
