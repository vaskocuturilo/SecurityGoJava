package com.example.java.mapper;


import com.example.java.dto.UserDto;
import com.example.java.entity.UserEntity;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring", unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface UserMapper {

    @Mapping(target = "active", source = "enabled")
    @Mapping(target = "roles", source = "roles")
    UserDto toUserDto(UserEntity user);
}
