(ns yocal-cljs.core
    (:require [reagent.core :as reagent]
              [re-frame.core :as re-frame]
              [yocal-cljs.handlers]
              [yocal-cljs.subs]
              [yocal-cljs.routes :as routes]
              [yocal-cljs.views :as views]))

(enable-console-print!)

(defn mount-root []
  (reagent/render [views/main-panel]
                  (.getElementById js/document "app")))

(defn ^:export init []
  (routes/app-routes)
  (re-frame/dispatch-sync [:initialize-db])
  (mount-root))
