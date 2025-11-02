#!/bin/bash

# Test script untuk Yard Planning API

BASE_URL="http://localhost:8080"

echo "🚢 Testing Yard Planning API"
echo "================================"

# Test 1: Health Check
echo -e "\n1️⃣  Testing Health Check..."
curl -s -X GET "$BASE_URL/health"
echo ""

# Test 2: Get Suggestion for 20ft container
echo -e "\n2️⃣  Testing Suggestion for 20ft DRY container..."
curl -s -X POST "$BASE_URL/suggestion" \
  -H "Content-Type: application/json" \
  -d '{
    "yard": "YRD1",
    "container_number": "ALFI000001",
    "container_size": 20,
    "container_height": 8.6,
    "container_type": "DRY"
  }' | jq
echo ""

# Test 3: Place Container
echo -e "\n3️⃣  Testing Placement..."
curl -s -X POST "$BASE_URL/placement" \
  -H "Content-Type: application/json" \
  -d '{
    "yard": "YRD1",
    "container_number": "ALFI000001",
    "block": "LC01",
    "slot": 1,
    "row": 1,
    "tier": 1
  }' | jq
echo ""

# Test 4: Get Suggestion for another container (should get slot 2)
echo -e "\n4️⃣  Testing Suggestion for second container..."
curl -s -X POST "$BASE_URL/suggestion" \
  -H "Content-Type: application/json" \
  -d '{
    "yard": "YRD1",
    "container_number": "ALFI000002",
    "container_size": 20,
    "container_height": 8.6,
    "container_type": "DRY"
  }' | jq
echo ""

# Test 5: Place second container
echo -e "\n5️⃣  Testing Placement of second container..."
curl -s -X POST "$BASE_URL/placement" \
  -H "Content-Type: application/json" \
  -d '{
    "yard": "YRD1",
    "container_number": "ALFI000002",
    "block": "LC01",
    "slot": 2,
    "row": 1,
    "tier": 1
  }' | jq
echo ""

# Test 6: Test stacking (tier 2 on same slot)
echo -e "\n6️⃣  Testing stacking on tier 2..."
curl -s -X POST "$BASE_URL/placement" \
  -H "Content-Type: application/json" \
  -d '{
    "yard": "YRD1",
    "container_number": "ALFI000003",
    "block": "LC01",
    "slot": 1,
    "row": 1,
    "tier": 2
  }' | jq
echo ""

# Test 7: Try to pickup blocked container (should fail)
echo -e "\n7️⃣  Testing pickup blocked container (should fail)..."
curl -s -X POST "$BASE_URL/pickup" \
  -H "Content-Type: application/json" \
  -d '{
    "yard": "YRD1",
    "container_number": "ALFI000001"
  }' | jq
echo ""

# Test 8: Pickup top container
echo -e "\n8️⃣  Testing pickup top container..."
curl -s -X POST "$BASE_URL/pickup" \
  -H "Content-Type: application/json" \
  -d '{
    "yard": "YRD1",
    "container_number": "ALFI000003"
  }' | jq
echo ""

# Test 9: Now pickup first container (should succeed)
echo -e "\n9️⃣  Testing pickup first container (should succeed now)..."
curl -s -X POST "$BASE_URL/pickup" \
  -H "Content-Type: application/json" \
  -d '{
    "yard": "YRD1",
    "container_number": "ALFI000001"
  }' | jq
echo ""

# Test 10: Get suggestion for 40ft container
echo -e "\n🔟 Testing Suggestion for 40ft container..."
curl -s -X POST "$BASE_URL/suggestion" \
  -H "Content-Type: application/json" \
  -d '{
    "yard": "YRD1",
    "container_number": "ALFI000004",
    "container_size": 40,
    "container_height": 8.6,
    "container_type": "DRY"
  }' | jq
echo ""

echo -e "\n✅ All tests completed!"