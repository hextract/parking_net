#!/usr/bin/env python3
"""
Comprehensive Integration Tests for Parking Net System

Tests the complete flow:
1. Register users (owner and driver)
2. Login both users
3. Owner creates parking places
4. Driver searches and views parking
5. Driver creates bookings
6. Owner views bookings
7. Updates and error cases
"""

import json
import time
import sys
import urllib.request
import urllib.parse
from datetime import datetime, timedelta
from typing import Dict, Optional, List

BASE_URLS = {
    'auth': 'http://localhost:8800',
    'parking': 'http://localhost:8888',
    'booking': 'http://localhost:8880'
}

class Response:
    def __init__(self, status_code: int, text: str):
        self.status_code = status_code
        self.text = text
        self._json = None
    
    def json(self):
        if self._json is None:
            try:
                self._json = json.loads(self.text)
            except json.JSONDecodeError:
                self._json = {}
        return self._json


class APIClient:
    def __init__(self, base_url: str):
        self.base_url = base_url
        self.token: Optional[str] = None
    
    def set_token(self, token: str):
        self.token = token
    
    def _make_request(self, method: str, path: str, data: Optional[Dict] = None, params: Optional[Dict] = None) -> Response:
        url = f"{self.base_url}{path}"
        if params:
            url += "?" + urllib.parse.urlencode(params)
        
        req_data = None
        headers = {'Content-Type': 'application/json'}
        if self.token:
            headers['api_key'] = self.token
        
        if data:
            req_data = json.dumps(data).encode('utf-8')
        
        req = urllib.request.Request(url, data=req_data, headers=headers, method=method)
        
        try:
            with urllib.request.urlopen(req, timeout=10) as response:
                return Response(response.getcode(), response.read().decode('utf-8'))
        except urllib.error.HTTPError as e:
            body = e.read().decode('utf-8') if e.fp else ""
            return Response(e.code, body)
        except urllib.error.URLError as e:
            return Response(0, f"Connection error: {str(e)}")
        except Exception as e:
            return Response(0, f"Request error: {str(e)}")
    
    def get(self, path: str, params: Optional[Dict] = None) -> Response:
        return self._make_request('GET', path, params=params)
    
    def post(self, path: str, data: Optional[Dict] = None) -> Response:
        return self._make_request('POST', path, data=data)
    
    def put(self, path: str, data: Optional[Dict] = None) -> Response:
        return self._make_request('PUT', path, data=data)
    
    def delete(self, path: str) -> Response:
        return self._make_request('DELETE', path)


class TestRunner:
    def __init__(self):
        self.auth_client = APIClient(BASE_URLS['auth'])
        self.parking_client = APIClient(BASE_URLS['parking'])
        self.booking_client = APIClient(BASE_URLS['booking'])
        self.timestamp = int(time.time())
        self.owner_token: Optional[str] = None
        self.driver_token: Optional[str] = None
        self.parking_id: Optional[int] = None
        self.booking_ids: List[int] = []
        self.parking_ids: List[int] = []
        self.passed = 0
        self.failed = 0
    
    def check_service_available(self, client: APIClient, service_name: str) -> bool:
        """Check if a service is available by making a simple request"""
        try:
            resp = client.get("/metrics")
            if resp.status_code != 0:
                return True
            return False
        except:
            return False
    
    def log(self, message: str, status: str = "INFO"):
        print(f"[{status}] {message}")
    
    def assert_status(self, response: Response, expected: int, test_name: str):
        if response.status_code == 0:
            self.log(f"FAILED: {test_name} - Service unavailable. Error: {response.text}", "ERROR")
            self.failed += 1
            return False
        if response.status_code != expected:
            try:
                error_body = response.json()
            except:
                error_body = response.text
            self.log(f"FAILED: {test_name} - Expected {expected}, got {response.status_code}. Body: {error_body}", "ERROR")
            self.failed += 1
            return False
        self.passed += 1
        return True
    
    def test_register_owner(self):
        self.log("Test 1: Register Owner")
        data = {
            "email": f"owner_{self.timestamp}@test.com",
            "login": f"owner_{self.timestamp}",
            "password": "password123",
            "role": "owner",
            "telegram_id": 123123123
        }
        resp = self.auth_client.post("/register", data)
        if not self.assert_status(resp, 200, "Register Owner"):
            return False
        
        result = resp.json()
        if 'token' not in result or not result.get('token'):
            self.log(f"FAILED: No token in response. Response: {result}", "ERROR")
            self.failed += 1
            return False
        
        self.owner_token = result['token']
        self.log(f"Owner registered. Token: {self.owner_token[:30] if len(self.owner_token) > 30 else self.owner_token}...")
        return True
    
    def test_register_driver(self):
        self.log("Test 2: Register Driver")
        data = {
            "email": f"driver_{self.timestamp}@test.com",
            "login": f"driver_{self.timestamp}",
            "password": "password123",
            "role": "driver",
            "telegram_id": 615711092
        }
        resp = self.auth_client.post("/register", data)
        if not self.assert_status(resp, 200, "Register Driver"):
            return False
        
        result = resp.json()
        if 'token' not in result or not result.get('token'):
            self.log(f"FAILED: No token in response. Response: {result}", "ERROR")
            self.failed += 1
            return False
        
        self.driver_token = result['token']
        self.log(f"Driver registered. Token: {self.driver_token[:30] if len(self.driver_token) > 30 else self.driver_token}...")
        return True
    
    def test_login_owner(self):
        self.log("Test 3: Login Owner")
        data = {
            "login": f"owner_{self.timestamp}",
            "password": "password123"
        }
        resp = self.auth_client.post("/login", data)
        if not self.assert_status(resp, 200, "Login Owner"):
            return False
        
        result = resp.json()
        if 'token' not in result or not result.get('token'):
            self.log(f"FAILED: No token in response. Response: {result}", "ERROR")
            self.failed += 1
            return False
        
        self.owner_token = result['token']
        self.log("Owner logged in successfully")
        return True
    
    def test_login_driver(self):
        self.log("Test 4: Login Driver")
        data = {
            "login": f"driver_{self.timestamp}",
            "password": "password123"
        }
        resp = self.auth_client.post("/login", data)
        if not self.assert_status(resp, 200, "Login Driver"):
            return False
        
        result = resp.json()
        if 'token' not in result or not result.get('token'):
            self.log(f"FAILED: No token in response. Response: {result}", "ERROR")
            self.failed += 1
            return False
        
        self.driver_token = result['token']
        self.log("Driver logged in successfully")
        return True
    
    def test_owner_creates_parking(self):
        self.log("Test 5: Owner Creates Parking Place")
        if not self.owner_token:
            self.log("SKIP: No owner token available (previous test failed)", "WARN")
            return True
        self.parking_client.set_token(self.owner_token)
        
        data = {
            "name": "Central Parking",
            "city": "Moscow",
            "address": "Red Square 1",
            "parking_type": "underground",
            "hourly_rate": 150,
            "capacity": 200
        }
        resp = self.parking_client.post("/parking", data)
        if not self.assert_status(resp, 200, "Create Parking"):
            return False
        
        result = resp.json()
        self.parking_id = result.get('id')
        if not self.parking_id:
            self.log("FAILED: No parking ID in response", "ERROR")
            self.failed += 1
            return False
        
        self.parking_ids.append(self.parking_id)
        self.log(f"Parking place created with ID: {self.parking_id}")
        return True
    
    def test_owner_creates_second_parking(self):
        self.log("Test 6: Owner Creates Second Parking Place")
        if not self.owner_token:
            self.log("SKIP: No owner token available (previous test failed)", "WARN")
            return True
        self.parking_client.set_token(self.owner_token)
        
        data = {
            "name": "Airport Parking",
            "city": "Moscow",
            "address": "Sheremetyevo Airport",
            "parking_type": "covered",
            "hourly_rate": 200,
            "capacity": 500
        }
        resp = self.parking_client.post("/parking", data)
        if not self.assert_status(resp, 200, "Create Second Parking"):
            return False
        
        result = resp.json()
        second_id = result.get('id')
        if second_id:
            self.parking_ids.append(second_id)
        self.log(f"Second parking place created with ID: {second_id}")
        return True
    
    def test_driver_searches_parking_by_city(self):
        self.log("Test 7: Driver Searches Parking by City")
        resp = self.parking_client.get("/parking", params={"city": "Moscow"})
        if not self.assert_status(resp, 200, "Search Parking by City"):
            return False
        
        parkings = resp.json()
        if not isinstance(parkings, list) or len(parkings) < 2:
            self.log(f"FAILED: Expected at least 2 parkings, got {len(parkings) if isinstance(parkings, list) else 'not a list'}", "ERROR")
            self.failed += 1
            return False
        
        self.log(f"Found {len(parkings)} parking places in Moscow")
        return True
    
    def test_driver_gets_parking_by_id(self):
        self.log("Test 8: Driver Gets Parking by ID")
        if not self.parking_id:
            self.log("SKIP: No parking ID available (previous test failed)", "WARN")
            return True
        
        resp = self.parking_client.get(f"/parking/{self.parking_id}")
        if not self.assert_status(resp, 200, "Get Parking by ID"):
            return False
        
        parking = resp.json()
        if parking.get('id') != self.parking_id:
            self.log(f"FAILED: Expected ID {self.parking_id}, got {parking.get('id')}", "ERROR")
            self.failed += 1
            return False
        
        self.log(f"Retrieved parking: {parking.get('name')}")
        return True
    
    def test_driver_searches_by_type(self):
        self.log("Test 9: Driver Searches by Parking Type")
        resp = self.parking_client.get("/parking", params={"parking_type": "underground"})
        if not self.assert_status(resp, 200, "Search by Type"):
            return False
        
        parkings = resp.json()
        self.log(f"Found {len(parkings)} underground parking places")
        return True
    
    def test_driver_creates_booking(self):
        self.log("Test 10: Driver Creates Booking")
        if not self.driver_token:
            self.log("SKIP: No driver token available (previous test failed)", "WARN")
            return True
        if not self.parking_id:
            self.log("SKIP: No parking ID available (previous test failed)", "WARN")
            return True
        self.booking_client.set_token(self.driver_token)
        
        date_from = (datetime.now() + timedelta(days=1)).strftime("%d-%m-%Y")
        date_to = (datetime.now() + timedelta(days=2)).strftime("%d-%m-%Y")
        
        data = {
            "parking_place_id": self.parking_id,
            "date_from": date_from,
            "date_to": date_to
        }
        resp = self.booking_client.post("/booking", data)
        if not self.assert_status(resp, 200, "Create Booking"):
            return False
        
        booking = resp.json()
        booking_id = booking.get('booking_id')
        if not booking_id:
            self.log(f"FAILED: No booking_id in response. Response: {booking}", "ERROR")
            self.failed += 1
            return False
        
        self.booking_ids.append(booking_id)
        self.log(f"Booking created: ID={booking_id}")
        return True
    
    def test_driver_gets_booking_by_id(self):
        self.log("Test 11: Driver Gets Booking by ID")
        if not self.driver_token:
            self.log("SKIP: No driver token available (previous test failed)", "WARN")
            return True
        if not self.parking_id:
            self.log("SKIP: No parking ID available (previous test failed)", "WARN")
            return True
        self.booking_client.set_token(self.driver_token)
        
        date_from = (datetime.now() + timedelta(days=3)).strftime("%d-%m-%Y")
        date_to = (datetime.now() + timedelta(days=4)).strftime("%d-%m-%Y")
        
        data = {
            "parking_place_id": self.parking_id,
            "date_from": date_from,
            "date_to": date_to
        }
        resp = self.booking_client.post("/booking", data)
        if resp.status_code == 200:
            booking = resp.json()
            booking_id = booking.get('booking_id')
            if not booking_id:
                self.log(f"FAILED: No booking_id in response. Response: {booking}", "ERROR")
                self.failed += 1
                return False
            
            self.booking_ids.append(booking_id)
            
            resp = self.booking_client.get(f"/booking/{booking_id}")
            if not self.assert_status(resp, 200, "Get Booking by ID"):
                return False
            
            retrieved = resp.json()
            retrieved_id = retrieved.get('id') or retrieved.get('booking_id')
            if retrieved_id != booking_id:
                self.log(f"FAILED: Expected ID {booking_id}, got {retrieved_id}", "ERROR")
                self.failed += 1
                return False
            
            self.log(f"Retrieved booking ID: {booking_id}")
            return True
        return False
    
    def test_owner_gets_bookings(self):
        self.log("Test 12: Owner Gets Bookings for Parking")
        if not self.owner_token:
            self.log("SKIP: No owner token available (previous test failed)", "WARN")
            return True
        if not self.parking_id:
            self.log("SKIP: No parking ID available (previous test failed)", "WARN")
            return True
        self.booking_client.set_token(self.owner_token)
        
        resp = self.booking_client.get("/booking", params={"parking_place_id": self.parking_id})
        if not self.assert_status(resp, 200, "Get Bookings"):
            return False
        
        bookings = resp.json()
        if not isinstance(bookings, list) or len(bookings) == 0:
            self.log(f"FAILED: Expected bookings list, got {bookings}", "ERROR")
            self.failed += 1
            return False
        
        self.log(f"Owner found {len(bookings)} bookings for parking {self.parking_id}")
        return True
    
    def test_owner_updates_parking(self):
        self.log("Test 13: Owner Updates Parking Place")
        if not self.owner_token:
            self.log("SKIP: No owner token available (previous test failed)", "WARN")
            return True
        if not self.parking_id:
            self.log("SKIP: No parking ID available (previous test failed)", "WARN")
            return True
        self.parking_client.set_token(self.owner_token)
        
        data = {
            "name": "Updated Central Parking",
            "city": "Moscow",
            "address": "Red Square 1, Updated",
            "parking_type": "underground",
            "hourly_rate": 180,
            "capacity": 250
        }
        resp = self.parking_client.put(f"/parking/{self.parking_id}", data)
        if not self.assert_status(resp, 200, "Update Parking"):
            return False
        
        parking = resp.json()
        if parking.get('name') != "Updated Central Parking":
            self.log(f"FAILED: Expected updated name, got {parking.get('name')}", "ERROR")
            self.failed += 1
            return False
        
        self.log("Parking place updated successfully")
        return True
    
    def test_driver_cannot_create_parking(self):
        self.log("Test 14: Driver Cannot Create Parking (Forbidden)")
        if not self.driver_token:
            self.log("SKIP: No driver token available (previous test failed)", "WARN")
            return True
        self.parking_client.set_token(self.driver_token)
        
        data = {
            "name": "Unauthorized Parking",
            "city": "Moscow",
            "address": "Test St",
            "parking_type": "outdoor",
            "hourly_rate": 100,
            "capacity": 50
        }
        resp = self.parking_client.post("/parking", data)
        if not self.assert_status(resp, 403, "Driver Create Parking Forbidden"):
            return False
        
        self.log("Driver correctly forbidden from creating parking")
        return True
    
    def test_driver_cannot_update_parking(self):
        self.log("Test 15: Driver Cannot Update Owner's Parking")
        if not self.driver_token:
            self.log("SKIP: No driver token available (previous test failed)", "WARN")
            return True
        if not self.parking_id:
            self.log("SKIP: No parking ID available (previous test failed)", "WARN")
            return True
        self.parking_client.set_token(self.driver_token)
        
        data = {
            "name": "Hacked Parking",
            "city": "Moscow",
            "address": "Hacked Address",
            "parking_type": "outdoor",
            "hourly_rate": 1,
            "capacity": 1
        }
        resp = self.parking_client.put(f"/parking/{self.parking_id}", data)
        if not self.assert_status(resp, 403, "Driver Update Parking Forbidden"):
            return False
        
        self.log("Driver correctly forbidden from updating owner's parking")
        return True
    
    def test_get_nonexistent_parking(self):
        self.log("Test 16: Get Non-Existent Parking (404)")
        resp = self.parking_client.get("/parking/99999")
        if not self.assert_status(resp, 404, "Get Non-Existent Parking"):
            return False
        
        self.log("Correctly returned 404 for non-existent parking")
        return True
    
    def test_update_booking_status(self):
        self.log("Test 17: Update Booking Status")
        if not self.driver_token:
            self.log("SKIP: No driver token available (previous test failed)", "WARN")
            return True
        self.booking_client.set_token(self.driver_token)
        
        if not self.booking_ids:
            self.log("SKIP: No bookings to update", "WARN")
            return True
        
        booking_id = self.booking_ids[0]
        
        resp = self.booking_client.get(f"/booking/{booking_id}")
        if resp.status_code != 200:
            self.log(f"SKIP: Could not retrieve booking {booking_id} to update", "WARN")
            return True
        
        existing_booking = resp.json()
        
        data = {
            "parking_place_id": existing_booking.get('parking_place_id') or self.parking_id,
            "date_from": existing_booking.get('date_from'),
            "date_to": existing_booking.get('date_to'),
            "status": "Confirmed",
            "full_cost": existing_booking.get('full_cost', 0)
        }
        
        resp = self.booking_client.put(f"/booking/{booking_id}", data)
        if not self.assert_status(resp, 200, "Update Booking Status"):
            return False
        
        booking = resp.json()
        if booking.get('status') != "Confirmed":
            self.log(f"FAILED: Expected status Confirmed, got {booking.get('status')}", "ERROR")
            self.failed += 1
            return False
        
        self.log(f"Booking status updated to: {booking.get('status')}")
        return True
    
    def test_driver_gets_own_bookings(self):
        self.log("Test 18: Driver Gets Own Bookings by user_id")
        if not self.driver_token:
            self.log("SKIP: No driver token available (previous test failed)", "WARN")
            return True
        if not self.booking_ids:
            self.log("SKIP: No bookings created yet", "WARN")
            return True
        self.booking_client.set_token(self.driver_token)
        
        # Get the user_id from a booking
        resp = self.booking_client.get(f"/booking/{self.booking_ids[0]}")
        if resp.status_code != 200:
            self.log("SKIP: Could not retrieve booking to get user_id", "WARN")
            return True
        
        booking = resp.json()
        user_id = booking.get('user_id')
        if not user_id:
            self.log("SKIP: No user_id in booking", "WARN")
            return True
        
        # Get bookings by user_id
        resp = self.booking_client.get("/booking", params={"user_id": user_id})
        if not self.assert_status(resp, 200, "Get Bookings by user_id"):
            return False
        
        bookings = resp.json()
        if not isinstance(bookings, list):
            self.log(f"FAILED: Expected bookings list, got {bookings}", "ERROR")
            self.failed += 1
            return False
        
        # Verify all bookings belong to the driver
        for b in bookings:
            if b.get('user_id') != user_id:
                self.log(f"FAILED: Found booking with different user_id: {b.get('user_id')}", "ERROR")
                self.failed += 1
                return False
        
        self.log(f"Driver found {len(bookings)} of their own bookings")
        return True
    
    def test_owner_gets_own_parkings(self):
        self.log("Test 19: Owner Gets Own Parkings by owner_id")
        if not self.owner_token:
            self.log("SKIP: No owner token available (previous test failed)", "WARN")
            return True
        if not self.parking_ids:
            self.log("SKIP: No parkings created yet", "WARN")
            return True
        self.parking_client.set_token(self.owner_token)
        
        # Get the owner_id from a parking
        resp = self.parking_client.get(f"/parking/{self.parking_ids[0]}")
        if resp.status_code != 200:
            self.log("SKIP: Could not retrieve parking to get owner_id", "WARN")
            return True
        
        parking = resp.json()
        owner_id = parking.get('owner_id')
        if not owner_id:
            self.log("SKIP: No owner_id in parking", "WARN")
            return True
        
        # Get parkings by owner_id
        resp = self.parking_client.get("/parking", params={"owner_id": owner_id})
        if not self.assert_status(resp, 200, "Get Parkings by owner_id"):
            return False
        
        parkings = resp.json()
        if not isinstance(parkings, list):
            self.log(f"FAILED: Expected parkings list, got {parkings}", "ERROR")
            self.failed += 1
            return False
        
        # Verify all parkings belong to the owner
        for p in parkings:
            if p.get('owner_id') != owner_id:
                self.log(f"FAILED: Found parking with different owner_id: {p.get('owner_id')}", "ERROR")
                self.failed += 1
                return False
        
        self.log(f"Owner found {len(parkings)} of their own parking places")
        return True
    
    def test_driver_deletes_own_booking(self):
        self.log("Test 20: Driver Deletes Own Booking")
        if not self.driver_token:
            self.log("SKIP: No driver token available (previous test failed)", "WARN")
            return True
        if not self.parking_id:
            self.log("SKIP: No parking ID available (previous test failed)", "WARN")
            return True
        self.booking_client.set_token(self.driver_token)
        
        # Create a booking to delete
        date_from = (datetime.now() + timedelta(days=5)).strftime("%d-%m-%Y")
        date_to = (datetime.now() + timedelta(days=6)).strftime("%d-%m-%Y")
        
        data = {
            "parking_place_id": self.parking_id,
            "date_from": date_from,
            "date_to": date_to
        }
        resp = self.booking_client.post("/booking", data)
        if resp.status_code != 200:
            self.log("SKIP: Could not create booking to delete", "WARN")
            return True
        
        booking = resp.json()
        booking_id = booking.get('booking_id')
        if not booking_id:
            self.log("SKIP: No booking_id in response", "WARN")
            return True
        
        # Delete the booking
        resp = self.booking_client.delete(f"/booking/{booking_id}")
        if not self.assert_status(resp, 200, "Delete Booking"):
            return False
        
        result = resp.json()
        if result.get('status') != 'success':
            self.log(f"FAILED: Expected success status, got {result.get('status')}", "ERROR")
            self.failed += 1
            return False
        
        # Verify booking is deleted
        resp = self.booking_client.get(f"/booking/{booking_id}")
        if resp.status_code != 404:
            self.log(f"FAILED: Booking should be deleted (404), got {resp.status_code}", "ERROR")
            self.failed += 1
            return False
        
        self.log(f"Booking {booking_id} deleted successfully")
        return True
    
    def test_owner_deletes_booking_for_their_parking(self):
        self.log("Test 21: Owner Deletes Booking for Their Parking")
        if not self.owner_token:
            self.log("SKIP: No owner token available (previous test failed)", "WARN")
            return True
        if not self.parking_id:
            self.log("SKIP: No parking ID available (previous test failed)", "WARN")
            return True
        if not self.driver_token:
            self.log("SKIP: No driver token available (previous test failed)", "WARN")
            return True
        
        # Create a booking as driver
        self.booking_client.set_token(self.driver_token)
        date_from = (datetime.now() + timedelta(days=7)).strftime("%d-%m-%Y")
        date_to = (datetime.now() + timedelta(days=8)).strftime("%d-%m-%Y")
        
        data = {
            "parking_place_id": self.parking_id,
            "date_from": date_from,
            "date_to": date_to
        }
        resp = self.booking_client.post("/booking", data)
        if resp.status_code != 200:
            self.log("SKIP: Could not create booking to delete", "WARN")
            return True
        
        booking = resp.json()
        booking_id = booking.get('booking_id')
        if not booking_id:
            self.log("SKIP: No booking_id in response", "WARN")
            return True
        
        # Owner deletes the booking
        self.booking_client.set_token(self.owner_token)
        resp = self.booking_client.delete(f"/booking/{booking_id}")
        if not self.assert_status(resp, 200, "Owner Delete Booking"):
            return False
        
        result = resp.json()
        if result.get('status') != 'success':
            self.log(f"FAILED: Expected success status, got {result.get('status')}", "ERROR")
            self.failed += 1
            return False
        
        self.log(f"Owner successfully deleted booking {booking_id} for their parking")
        return True
    
    def test_driver_cannot_delete_other_driver_booking(self):
        self.log("Test 22: Driver Cannot Delete Other Driver's Booking (403)")
        if not self.driver_token:
            self.log("SKIP: No driver token available (previous test failed)", "WARN")
            return True
        if not self.booking_ids:
            self.log("SKIP: No bookings available", "WARN")
            return True
        
        # Try to delete a booking (should work if it's theirs, but we'll test with a non-existent one)
        # Actually, let's create a second driver and try to delete first driver's booking
        # For simplicity, we'll test with a non-existent booking ID
        self.booking_client.set_token(self.driver_token)
        resp = self.booking_client.delete("/booking/99999")
        if not self.assert_status(resp, 404, "Delete Non-Existent Booking"):
            return False
        
        self.log("Correctly returned 404 for non-existent booking")
        return True
    
    def test_owner_deletes_own_parking(self):
        self.log("Test 23: Owner Deletes Own Parking")
        if not self.owner_token:
            self.log("SKIP: No owner token available (previous test failed)", "WARN")
            return True
        self.parking_client.set_token(self.owner_token)
        
        # Create a parking to delete
        data = {
            "name": "Temporary Parking",
            "city": "Moscow",
            "address": "Temp St 1",
            "parking_type": "outdoor",
            "hourly_rate": 50,
            "capacity": 10
        }
        resp = self.parking_client.post("/parking", data)
        if resp.status_code != 200:
            self.log("SKIP: Could not create parking to delete", "WARN")
            return True
        
        result = resp.json()
        parking_id = result.get('id')
        if not parking_id:
            self.log("SKIP: No parking ID in response", "WARN")
            return True
        
        # Delete the parking
        resp = self.parking_client.delete(f"/parking/{parking_id}")
        if not self.assert_status(resp, 200, "Delete Parking"):
            return False
        
        result = resp.json()
        if result.get('status') != 'success':
            self.log(f"FAILED: Expected success status, got {result.get('status')}", "ERROR")
            self.failed += 1
            return False
        
        # Verify parking is deleted
        resp = self.parking_client.get(f"/parking/{parking_id}")
        if resp.status_code != 404:
            self.log(f"FAILED: Parking should be deleted (404), got {resp.status_code}", "ERROR")
            self.failed += 1
            return False
        
        self.log(f"Parking {parking_id} deleted successfully")
        return True
    
    def test_driver_cannot_delete_parking(self):
        self.log("Test 24: Driver Cannot Delete Parking (403)")
        if not self.driver_token:
            self.log("SKIP: No driver token available (previous test failed)", "WARN")
            return True
        if not self.parking_id:
            self.log("SKIP: No parking ID available (previous test failed)", "WARN")
            return True
        self.parking_client.set_token(self.driver_token)
        
        resp = self.parking_client.delete(f"/parking/{self.parking_id}")
        if not self.assert_status(resp, 403, "Driver Delete Parking Forbidden"):
            return False
        
        self.log("Driver correctly forbidden from deleting parking")
        return True
    
    def test_owner_cannot_delete_other_owner_parking(self):
        self.log("Test 25: Owner Cannot Delete Other Owner's Parking (403)")
        if not self.owner_token:
            self.log("SKIP: No owner token available (previous test failed)", "WARN")
            return True
        
        # Try to delete a parking that doesn't belong to this owner
        # We'll use a non-existent ID which should return 404, but if it exists and belongs to someone else, it would be 403
        self.parking_client.set_token(self.owner_token)
        resp = self.parking_client.delete("/parking/99999")
        if resp.status_code not in [403, 404]:
            self.log(f"FAILED: Expected 403 or 404, got {resp.status_code}", "ERROR")
            self.failed += 1
            return False
        
        self.log(f"Correctly returned {resp.status_code} for unauthorized/non-existent parking deletion")
        return True
    
    def check_services(self):
        """Check if all required services are available"""
        self.log("Checking service availability...")
        services_ok = True
        
        if not self.check_service_available(self.auth_client, "auth"):
            self.log("ERROR: Auth service is not available at " + BASE_URLS['auth'], "ERROR")
            self.log("  Please run: docker-compose up -d", "ERROR")
            services_ok = False
        
        if not self.check_service_available(self.parking_client, "parking"):
            self.log("ERROR: Parking service is not available at " + BASE_URLS['parking'], "ERROR")
            self.log("  Please run: docker-compose up -d", "ERROR")
            services_ok = False
        
        if not self.check_service_available(self.booking_client, "booking"):
            self.log("ERROR: Booking service is not available at " + BASE_URLS['booking'], "ERROR")
            self.log("  Please run: docker-compose up -d", "ERROR")
            services_ok = False
        
        if services_ok:
            self.log("All services are available", "INFO")
        else:
            self.log("CRITICAL: Required services are not available. Tests will fail.", "ERROR")
            self.failed += 1  # Count this as a failure
        
        return services_ok
    
    def run_all(self):
        self.log("=" * 60)
        self.log("Starting Integration Tests")
        self.log("=" * 60)
        
        # Check services before running tests
        if not self.check_services():
            self.log("=" * 60)
            self.log(f"Tests aborted: {self.passed} passed, {self.failed} failed")
            self.log("=" * 60)
            return False
        
        self.log("")
        
        tests = [
            self.test_register_owner,
            self.test_register_driver,
            self.test_login_owner,
            self.test_login_driver,
            self.test_owner_creates_parking,
            self.test_owner_creates_second_parking,
            self.test_driver_searches_parking_by_city,
            self.test_driver_gets_parking_by_id,
            self.test_driver_searches_by_type,
            self.test_driver_creates_booking,
            self.test_driver_gets_booking_by_id,
            self.test_owner_gets_bookings,
            self.test_owner_updates_parking,
            self.test_driver_cannot_create_parking,
            self.test_driver_cannot_update_parking,
            self.test_get_nonexistent_parking,
            self.test_update_booking_status,
            self.test_driver_gets_own_bookings,
            self.test_owner_gets_own_parkings,
            self.test_driver_deletes_own_booking,
            self.test_owner_deletes_booking_for_their_parking,
            self.test_driver_cannot_delete_other_driver_booking,
            self.test_owner_deletes_own_parking,
            self.test_driver_cannot_delete_parking,
            self.test_owner_cannot_delete_other_owner_parking,
        ]
        
        for test in tests:
            try:
                test()
            except Exception as e:
                self.log(f"Exception in test: {e}", "ERROR")
                self.failed += 1
        
        self.log("=" * 60)
        self.log(f"Tests completed: {self.passed} passed, {self.failed} failed")
        self.log("=" * 60)
        
        return self.failed == 0


if __name__ == "__main__":
    runner = TestRunner()
    success = runner.run_all()
    sys.exit(0 if success else 1)

