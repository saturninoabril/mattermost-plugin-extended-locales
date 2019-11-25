
const {addLocaleData} = window.ReactIntl;

// For demo purpose only
const noLocaleData = require('react-intl/locale-data/no');
const tlLocaleData = require('react-intl/locale-data/tl');

const fiLocaleData = require('react-intl/locale-data/fi');
const idLocaleData = require('react-intl/locale-data/id');
const seLocaleData = require('react-intl/locale-data/se');

import {getLocales} from './actions';
import {id as pluginId} from './manifest';

export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        console.log('initialize ExtendedLocalesPlugin')
        store.dispatch(getLocales());

        registry.registerWebSocketEventHandler(
            'custom_' + pluginId + '_locales_change',
            () => {store.dispatch(getLocales());},
        );

        // Manually add locales
        // (future option could be to add registry on webapp to register adding locales)
        addLocaleData(noLocaleData);
        addLocaleData(tlLocaleData);

        addLocaleData(fiLocaleData);
        addLocaleData(idLocaleData);
        addLocaleData(seLocaleData);
    }
}

window.registerPlugin(pluginId, new Plugin());
