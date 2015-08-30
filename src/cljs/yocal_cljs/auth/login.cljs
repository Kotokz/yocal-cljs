(ns yocal-cljs.auth.login
  (:require [reagent.core :as ratom]
            [reforms.reagent :include-macros true :as f]
            [reforms.validation :as v :include-macros true]))

(def auth-validators
  [(v/present [:username] "Enter login name")
   (v/present [:password] "Enter password")])
(def data (ratom/atom {}))

(defn sign-in!
 [data credentials]
 (when (apply v/validate! data credentials auth-validators)
   #(pr credentials)))

(defn login-form-view
  [credentials]
    (v/form credentials
      (f/group-title {:class "group-title-main" :key "title"} "Please login with your user name")
      (v/text {:key "login"} "Login" data [:username] :placeholder "Enter login")
      (v/password {:key "pwd"} "Password" data [:password] :placeholder "Enter password")
      (f/form-buttons
        (f/button {:key "save"} "Save" #(sign-in! data credentials)))))


(f/set-options!
  {:form        {:horizontal true
                 :attrs      {:style {:border           "1px solid #BBB"
                                      :border-radius    "5px"
                                      :background-color "#EFEFEF"
                                      :padding          "10px 20px"
                                      :width            "550px"}}}
   :group-title {:tag   :h3
                 :attrs {:style {:color "#666"}}}})
