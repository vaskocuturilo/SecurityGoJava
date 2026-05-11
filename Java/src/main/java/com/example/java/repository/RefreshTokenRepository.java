package com.example.java.repository;


import com.example.java.entity.RefreshTokenEntity;
import com.example.java.entity.UserEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface RefreshTokenRepository extends JpaRepository<RefreshTokenEntity, String> {
    Optional<RefreshTokenEntity> findByToken(String token);

    @Modifying
    @Query("UPDATE RefreshTokenEntity t SET t.revoked = true WHERE t.user = :user")
    List<RefreshTokenEntity> findAllByUser(UserEntity user);

    void deleteByUser(UserEntity user);
}
