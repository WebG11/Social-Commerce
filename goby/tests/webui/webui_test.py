from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import unittest

class WebUITest(unittest.TestCase):
    def setUp(self):
        self.driver = webdriver.Chrome()

    def test_home_page_title(self):
        self.driver.get("http://localhost:3000")
        self.assertEqual(self.driver.title, "GoBy")

    def test_login_functionality(self):
        self.driver.get("http://localhost:3000/seller/login")
        
        # 添加等待确保页面加载完成
        wait = WebDriverWait(self.driver, 10)
        
        # 使用ID来查找元素
        email_field = wait.until(EC.presence_of_element_located((By.ID, "email")))
        email_field.send_keys("22009201315@stu.xidian.edu.cn")
        
        # 密码和按钮字段也需要检查实际的HTML属性
        # 你需要检查密码输入框和提交按钮的实际id或其他属性
        password_field = self.driver.find_element(By.ID, "password")  # 假设密码框id是"password"
        password_field.send_keys("123456")
        
        # 查找提交按钮，通过按钮文本定位
        submit_button = wait.until(EC.element_to_be_clickable((By.XPATH, "//button[text()='Continue']")))
        submit_button.click()
        
        #  检查URL是否改变
        self.assertNotEqual(self.driver.current_url, "http://localhost:3000")

    def tearDown(self):
        self.driver.quit()

if __name__ == "__main__":
    unittest.main()