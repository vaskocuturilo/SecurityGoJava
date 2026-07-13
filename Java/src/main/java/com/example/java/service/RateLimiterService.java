package com.example.java.service;


import io.github.bucket4j.Bandwidth;
import io.github.bucket4j.Bucket;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.time.Duration;
import java.util.concurrent.ConcurrentHashMap;

@Slf4j
@Service
@RequiredArgsConstructor
public class RateLimiterService {
    private final ConcurrentHashMap<String, Bucket> buckets = new ConcurrentHashMap<>();

    public boolean isAllowed(final String key, int limit, int durationSeconds) {
        final Bucket bucket = buckets.computeIfAbsent(key,
                k -> createBucket(limit, durationSeconds));

        final boolean allowed = bucket.tryConsume(1);

        if (!allowed) {
            log.warn("Rate limit exceeded for key [{}]: {} requests per {}s",
                    key, limit, durationSeconds);
        }

        return allowed;
    }

    private Bucket createBucket(int limit, int durationSeconds) {
        final Bandwidth bandwidth = Bandwidth.builder()
                .capacity(limit)
                .refillGreedy(limit, Duration.ofSeconds(durationSeconds))
                .build();

        return Bucket.builder()
                .addLimit(bandwidth)
                .build();
    }
}
