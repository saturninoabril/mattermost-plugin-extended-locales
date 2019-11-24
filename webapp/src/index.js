
const {addLocaleData} = window.ReactIntl;

// import '@formatjs/intl-relativetimeformat/polyfill';
// import '@formatjs/intl-relativetimeformat/dist/locale-data/en';
// import '@formatjs/intl-relativetimeformat/dist/locale-data/fr'

// For demo purpose only
const noLocaleData = require('react-intl/locale-data/no');
const tlLocaleData = require('react-intl/locale-data/tl');

import {getLocales} from './actions';

import {id as pluginId} from './manifest';

export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        console.log('initialize ExtendedLocalesPlugin')
        store.dispatch(getLocales());

        // Manually add locales
        // (future option could be to add registry on webapp to register adding locales)
        addLocaleData(noLocaleData);
        addLocaleData(tlLocaleData);
    }
}

window.registerPlugin(pluginId, new Plugin());
