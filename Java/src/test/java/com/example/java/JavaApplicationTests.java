package com.example.java;

import com.example.java.rest.RestTaskControllerV1;
import com.example.java.rest.UserRestControllerV1;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;

import static org.assertj.core.api.AssertionsForClassTypes.assertThat;

@ActiveProfiles("test")
@SpringBootTest
class JavaApplicationTests extends AbstractRestControllerBaseTest {

    @Autowired
    private RestTaskControllerV1 restTaskControllerV1;

    @Autowired
    private UserRestControllerV1 userRestControllerV1;

    @Test
    void contextLoads() {
        assertThat(restTaskControllerV1).isNotNull();
        assertThat(userRestControllerV1).isNotNull();
    }
}
