"""
A simple selenium test example written by python
"""

import unittest
import binascii
from selenium import webdriver
from selenium.common.exceptions import NoSuchElementException
from selenium.webdriver.common.keys import Keys
import subprocess
import sys
import time

def jiggleMouse(driver, sendKeyTarget):
    if sendKeyTarget != "":
        try:
            driver.find_element_by_class_name(sendKeyTarget).send_keys("s")
        except:
            pass

def killVlc(dockerId):
     # find and kill vlc ########### Don't Change
    pid = subprocess.check_output(["docker", "exec", dockerId, "pidof", "vlc"]).decode("ascii").strip()
    subprocess.check_output(["docker", "exec", dockerId, "kill", "-INT", pid])

def process():
    global wdriver
    # config for the specific service
    firstPlaybackButton = sys.argv[1]
    errorClass = sys.argv[2]
    titleClass = sys.argv[3]
    fullscreenButton = sys.argv[4]
    seccondPlaybackButton = sys.argv[5]
    fp = webdriver.FirefoxProfile(sys.argv[6])
    pauseButton = sys.argv[7]
    url = sys.argv[8]
    dockerId = sys.argv[9]
    maxCheckLoop = int(sys.argv[10])
    sendKeyTarget = sys.argv[11]
    finishedClass = sys.argv[12]
    # end config for the specific service

    # start the browser
    options = webdriver.FirefoxOptions()
    options.profile = fp
    wdriver = webdriver.Remote(command_executor="http://localhost:4444/wd/hub", desired_capabilities=options.to_capabilities())
    
    # start url
    wdriver.get(url)

    # click the play button, not always necesary which is why the loop can end without clicking it.
    if firstPlaybackButton != "":
        done = False
        error = None
        for i in range(maxCheckLoop):
            try:
                jiggleMouse(wdriver, sendKeyTarget)
                firstPlayButton = wdriver.find_element_by_class_name(firstPlaybackButton).click()
                done = True
                break
            except Exception as e:
                error = e
                time.sleep(1)
        if not done:
            killVlc(dockerId)
            raise error

    # load the title
    episode = ""
    error = None
    for i in range(maxCheckLoop):
        try:
            jiggleMouse(wdriver, sendKeyTarget)
            episode = wdriver.find_element_by_class_name(titleClass).get_attribute('innerHTML')
            if episode != "":
                break
        except Exception as e:
            error = e
        time.sleep(1)
    if episode == "":
        raise error

    # check if it is a subed show
    if "(Sub)" in episode:
        print("this is a subtitled show")
        wdriver.quit()
        sys.exit(2)
    
    # load error
    if errorClass != "":
        for i in range(maxCheckLoop):
            try:
                jiggleMouse(wdriver, sendKeyTarget)
                error = wdriver.find_element_by_class_name(errorClass).get_attribute('innerHTML')
                break
            except:
                pass
            time.sleep(1)
        
        # print the error
        if error != "" :
            print(error)
            killVlc(dockerId)
            wdriver.quit()
            sys.exit(1)
    
    # print the title
    print(episode)

    # click the fullscreen button
    done = False
    error = None
    for i in range(maxCheckLoop):
        try:
            jiggleMouse(wdriver, sendKeyTarget)
            full = wdriver.find_element_by_class_name(fullscreenButton).click()
            done = True
            break
        except Exception as e:
            error = e
            time.sleep(1)
    if not done:
        killVlc(dockerId)
        raise error


    if seccondPlaybackButton != "":
        done = False
        error = None
        for i in range(maxCheckLoop):
            try:
                jiggleMouse(wdriver, sendKeyTarget)
                try:
                    if wdriver.find_element_by_class_name(pauseButton).is_displayed():
                        done = True
                        break
                except:
                    pass
                wdriver.find_element_by_class_name(seccondPlaybackButton).click()
                done = True
                break
            except Exception as e:
                error = e
                time.sleep(1)
        if not done:
            killVlc(dockerId)
            raise error

    killVlc(dockerId)

    # wait for next episode to end
    newEpisode = episode
    testUrl = wdriver.current_url
    while episode == newEpisode:
        time.sleep(10)
        if testUrl != wdriver.current_url:
            break
        if sendKeyTarget != "" and (len(wdriver.find_elements_by_class_name(sendKeyTarget)) ==  0 or not wdriver.find_elements_by_class_name(sendKeyTarget)[0].is_displayed()):
            break
        if finishedClass != "" and len(wdriver.find_elements_by_class_name(finishedClass)) > 0:
            break
        try :
            newEpisode = wdriver.find_element_by_class_name(titleClass).get_attribute('innerHTML')
        except:
            pass

    # kill ffmpeg
    pid = subprocess.check_output(["docker", "exec", dockerId, "pidof", "ffmpeg"]).decode("ascii").strip()
    subprocess.check_output(["docker", "exec", dockerId, "kill", "-INT", pid])
    wdriver.quit()

if __name__ == '__main__':
    global wdriver
    try:
        process()
    except Exception as e:
        try:
            dockerid = sys.argv[9]
            killVlc(sys.argv[9])
        except:
            pass
        try:
            pid = subprocess.check_output(["docker", "exec", dockerId, "pidof", "ffmpeg"]).decode("ascii").strip()
            if pid != "":
                time.sleep(5)
                subprocess.check_output(["docker", "exec", dockerId, "kill", "-INT", pid])
        except:
            pass
        try:
            wdriver.quit()
        except Exception as b:
            print(b)
            print()
        raise e