(ns yocal-cljs.view.signup
  (:require [reagent.core :as ratom]
            [clojure.string :as str]
            [re-frame.core :as re-frame]
            [reagent-forms.core :refer [bind-fields init-field value-of]]
            [yocal-cljs.view.shared :refer [row input]]
            [secretary.core :as secretary]))

(def signup-form-data (ratom/atom {:user       {:username ""
                                                :fullname ""
                                                :staffid  nil
                                                :email    ""
                                                :password ""
                                                :retype   ""}
                                   :is-loading false}))
(defn sign-up! []
  (let [name (get-in @signup-form-data [:user :username])
        fname (get-in @signup-form-data [:user :fullname])
        sid (get-in @signup-form-data [:user :staffid])
        email (get-in @signup-form-data [:user :email])
        pwd (get-in @signup-form-data [:user :password])
        retype (get-in @signup-form-data [:user :retype])]
    (cond
      (empty? name) (swap! signup-form-data assoc-in [:errors :username] "User name is empty")
      (empty? fname) (swap! signup-form-data assoc-in [:errors :fullname] "Full name is empty")
      (not (= 8 (count (str sid)))) (swap! signup-form-data assoc-in [:errors :staffid] "StaffID can only be a 8 digits number")
      (empty? email) (swap! signup-form-data assoc-in [:errors :email] "Email is empty")
      (empty? pwd) (swap! signup-form-data assoc-in [:errors :password] "Password is empty")
      (empty? retype) (swap! signup-form-data assoc-in [:errors :retype] "Password is empty")
      (not (= pwd retype)) (swap! signup-form-data assoc-in [:errors :retype] "Please put same password in both fields")
      :else (re-frame/dispatch [:sign-up (:user @signup-form-data)]))))


(def form-login
  [:div
   (input "User Name" :text :user.username)
   [:div.row
    [:div.col-md-2]
    [:div.col-md-5
     [:div.alert.alert-danger
      {:field :alert :id :errors.username}]]]

   (input "Full Name" :text :user.fullname)
   [:div.row
    [:div.col-md-2]
    [:div.col-md-5
     [:div.alert.alert-danger
      {:field :alert :id :errors.fullname}]]]

   (input "Staff ID" :numeric :user.staffid)
   [:div.row
    [:div.col-md-2]
    [:div.col-md-5
     [:div.alert.alert-danger
      {:field :alert :id :errors.staffid}]]]

   (input "Email" :email :user.email)
   [:div.row
    [:div.col-md-2]
    [:div.col-md-5
     [:div.alert.alert-danger
      {:field :alert :id :errors.email}]]]

   (input "password" :password :user.password)
   [:div.row
    [:div.col-md-2]
    [:div.col-md-5
     [:div.alert.alert-danger
      {:field :alert :id :errors.password}]]]

   (input "password" :password :user.retype)
   [:div.row
    [:div.col-md-2]
    [:div.col-md-5
     [:div.alert.alert-danger
      {:field :alert :id :errors.retype}]
     [:div.alert.alert-danger
      {:field :alert :id :errors.other}]]]])



(def signup-form-page
  (fn []
    (let [isloading (:is-loading @signup-form-data)
          jwt (:jwt @(re-frame/subscribe [:user]))]
      (if-not (str/blank? jwt) (secretary/dispatch! "/home") (pr "blank jwt"))
      [:div.col-md-6
       [:div.padding]
       [:div.page-header [:h1 "Registration Form"]]
       [bind-fields
        form-login
        signup-form-data]
       [:button.btn.btn-default {:on-click #(sign-up!)}
        (if isloading [:span.glyphicon.glyphicon-refresh.glyphicon-refresh-animate])
        "Sign Up!"]])))

