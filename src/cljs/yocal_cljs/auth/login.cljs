(ns yocal-cljs.auth.login
  (:require [reagent.core :as ratom]
            [re-frame.core :as re-frame]
            [reforms.reagent :include-macros true :as f]
            [reforms.validation :as v :include-macros true]))

(def auth-validators
  [(v/present [:username] "Enter login name")
   (v/present [:password] "Enter password")])

(def data (ratom/atom {}))
(def credentials (ratom/atom {:username "" :password ""}))

(defn sign-in!
 [data credentials]
 (when (apply v/validate! data credentials auth-validators)
   (re-frame/dispatch [:login @data])))

(defn login-form-view
  []
   (let [jwt @(re-frame/subscribe [:auth-jwt])]
    (v/form credentials {:key "form"}
      (f/group-title {:class "group-title-main" :key "title"} "Please login with your user name")
      (v/text {:key "login"} "Login" data [:username] :placeholder "Enter login")
      (v/password {:key "pwd"} "Password" data [:password] :placeholder "Enter password")
      (f/form-buttons
        (f/button {:key "save"} "Save" #(sign-in! data credentials))
        (f/button {:key "get"} "Get" #(re-frame/dispatch [:get-balance jwt]))))))


(f/set-options!
  {:form        {:horizontal true
                 :attrs      {:style {:border           "1px solid #BBB"
                                      :border-radius    "5px"
                                      :background-color "#EFEFEF"
                                      :padding          "10px 20px"
                                      :width            "550px"}}}
   :group-title {:tag   :h3
                 :attrs {:style {:color "#666"}}}})
