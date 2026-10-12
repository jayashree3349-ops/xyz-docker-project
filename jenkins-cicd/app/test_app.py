import unittest
from app import Handler

class AppSourceTest(unittest.TestCase):
    def test_handler_has_get(self):
        self.assertTrue(callable(getattr(Handler, "do_GET", None)))

if __name__ == "__main__":
    unittest.main()
